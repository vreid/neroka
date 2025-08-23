import { createMistral } from "@ai-sdk/mistral";
import { createOpenAICompatible } from "@ai-sdk/openai-compatible";
import { generateObject } from "ai";
import { createRestAPIClient } from "masto";
import type { Status } from "masto/mastodon/entities/v1/status.js";
import { EOL } from "node:os";
import yargs from "yargs";
import { z } from "zod";
import { convertToDialog } from "./util";

import rewriter from "./rewriter.txt" with { type: "text" };

const args = yargs(process.argv)
  .help("help")
  .alias("h", "help")
  .version("version", "xyz")
  .alias("v", "version")
  .options({
    server: {
      alias: "s",
      require: true,
      type: "string",
      default: Bun.env["MASTODON_SERVER"],
    },
    "access-token": {
      alias: "t",
      require: true,
      type: "string",
      default: Bun.env["MASTODON_ACCESS_TOKEN"],
    },
    "mistral-api-key": {
      require: true,
      type: "string",
      default: Bun.env["MISTRAL_API_KEY"],
    },
    "mistral-base-url": {
      type: "string",
      default: Bun.env["MISTRAL_BASE_URL"],
    },
    "openrouter-api-key": {
      require: true,
      type: "string",
      default: Bun.env["OPENROUTER_API_KEY"],
    },
    "openrouter-base-url": {
      type: "string",
      default: Bun.env["OPENROUTER_BASE_URL"],
    },
  })
  .parseSync();

const mistral = createMistral({
  baseURL: args.mistralBaseUrl,
  apiKey: args.mistralApiKey,
});

const openrouter = createOpenAICompatible({
  name: "openrouter",
  baseURL: args.openrouterBaseUrl || "https://openrouter.ai/api/v1",
  apiKey: args.openrouterApiKey,
});

const restClient = createRestAPIClient({
  url: args.server,
  accessToken: args.accessToken,
});

const conversations = await restClient.v1.conversations.list({ limit: 100 });
const fullConversations = await Promise.all(
  conversations
    .filter((x) => x != null && x.lastStatus != null)
    .map(async (conversation) => {
      const lastStatus = conversation.lastStatus;
      if (lastStatus == null) {
        return {
          id: conversation.id,
          account: conversation.accounts[0]?.username || "",
          history: [] as Status[],
        };
      }

      const statusContext = await restClient.v1.statuses
        .$select(lastStatus.id)
        .context.fetch();

      const history = [
        ...statusContext.ancestors,
        lastStatus,
        ...statusContext.descendants,
      ];

      history.sort((x, y) => {
        if (x.createdAt > y.createdAt) {
          return 1;
        }

        if (x.createdAt < y.createdAt) {
          return -1;
        }

        return 0;
      });

      return {
        id: conversation.id,
        account: conversation.accounts[0]?.username || "",
        history,
      };
    })
);

for (const conversation of fullConversations) {
  const history = conversation.history;

  const dialog = await convertToDialog(mistral("pixtral-12b-latest"), history);

  const { object } = await generateObject({
    model: openrouter("google/gemini-2.5-pro"),
    schema: z.array(
      z.object({
        username: z.string(),
        message: z.string(),
      })
    ),
    prompt: [
      {
        role: "system",
        content: rewriter,
      },
      {
        role: "user",
        content: dialog.join(EOL),
      },
    ],
  });

  console.log(JSON.stringify(object));
}
