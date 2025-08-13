import { createAnthropic } from "@ai-sdk/anthropic";
import { createDeepSeek } from "@ai-sdk/deepseek";
import { createMistral } from "@ai-sdk/mistral";
import { createOpenAICompatible } from "@ai-sdk/openai-compatible";
import {
  generateText,
  type AssistantModelMessage,
  type LanguageModel,
  type UserModelMessage,
} from "ai";
import { convert } from "html-to-text";
import { createRestAPIClient, createStreamingAPIClient } from "masto";
import type { MediaAttachment } from "masto/mastodon/entities/v1/media-attachment.js";
import moment from "moment";
import { EOL } from "node:os";
import sharp from "sharp";
import yargs from "yargs";
import { hideBin } from "yargs/helpers";

// @ts-ignore
import neroka from "./neroka.md" with { type: "text" };

const args = yargs(hideBin(process.argv))
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
    "anthropic-api-key": {
      require: true,
      type: "string",
      default: Bun.env["ANTHROPIC_API_KEY"],
    },
    "anthropic-base-url": {
      type: "string",
      default: Bun.env["ANTHROPIC_BASE_URL"],
    },
    "deepseek-api-key": {
      require: true,
      type: "string",
      default: Bun.env["DEEPSEEK_API_KEY"],
    },
    "deepseek-base-url": {
      type: "string",
      default: Bun.env["DEEPSEEK_BASE_URL"],
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

const anthropic = createAnthropic({
  baseURL: args.anthropicBaseUrl,
  apiKey: args.anthropicApiKey,
});

const deepSeek = createDeepSeek({
  baseURL: args.deepseekBaseUrl,
  apiKey: args.deepseekApiKey,
});

const mistral = createMistral({
  baseURL: args.mistralBaseUrl,
  apiKey: args.mistralApiKey,
});

const openaiCompatible = createOpenAICompatible({
  name: "openrouter",
  baseURL: args.openrouterBaseUrl || "https://openrouter.ai/api/v1",
  apiKey: args.openrouterApiKey,
});

const restClient = createRestAPIClient({
  url: args.server,
  accessToken: args.accessToken,
});

const instance = await restClient.v2.instance.fetch();
const streamingClient = createStreamingAPIClient({
  streamingApiUrl: instance.configuration.urls.streaming,
  accessToken: args.accessToken,
});

const credentials = await restClient.v1.accounts.verifyCredentials();

function textContent(content: string): string {
  return convert(content, {
    wordwrap: false,
    selectors: [{ selector: "a", format: "skip" }],
  });
}

async function describeAttachment(mediaAttachment: MediaAttachment) {
  if (mediaAttachment.type != "image" || mediaAttachment.url == null) {
    return "";
  }

  console.log(`describing ${mediaAttachment.id}`);

  const response = await fetch(mediaAttachment.url);
  const buffer = await sharp(await response.arrayBuffer())
    .resize(1000)
    .webp()
    .toBuffer();

  const { text } = await generateText({
    model: mistral("pixtral-12b-latest"),
    messages: [
      {
        role: "user",
        content: [
          {
            type: "text",
            text: "Describe the image.",
          },
          {
            type: "image",
            image: buffer,
          },
        ],
      },
    ],
  });

  console.log(
    `description is ${text.substring(0, Math.min(100, text.length))}`
  );
  return text;
}

console.log("waiting for events");
for await (const event of streamingClient.direct.subscribe()) {
  console.log(`received ${event.event}`);
  if (event.event != "conversation") {
    continue;
  }

  const lastStatus = event.payload.lastStatus;
  if (!lastStatus) {
    console.log("lastStatus was null");
    continue;
  }

  const statusId = lastStatus.id;
  console.log(`id is ${statusId}`);

  const status = await restClient.v1.statuses.$select(statusId).fetch();
  if (status.account.username == credentials.username) {
    console.log("last post is from myself");
    continue;
  }

  const statusContext = await restClient.v1.statuses
    .$select(statusId)
    .context.fetch();

  console.log(
    `context has ${statusContext.ancestors.length} ancestors and ${statusContext.descendants.length} descendants`
  );

  const currentConversation = [
    ...statusContext.ancestors,
    status,
    ...statusContext.descendants,
  ];

  const userMessages = await Promise.all(
    currentConversation.map(async (entry) => {
      const username = entry.account.username;

      const createdAt = moment(entry.createdAt).format(
        "MMMM Do YYYY, h:mm:ss a"
      );

      const imageDescriptions = (
        await Promise.all(
          entry.mediaAttachments.map(
            async (mediaAttachment) => await describeAttachment(mediaAttachment)
          )
        )
      )
        .filter((x) => x && x.length > 0)
        .join(EOL);

      const content = (() => {
        if (imageDescriptions.length == 0) {
          return `${textContent(entry.content)}`;
        }

        return [`${textContent(entry.content)}`, imageDescriptions].join(EOL);
      })();

      if (username == credentials.username) {
        return {
          role: "assistant",
          content,
        } as AssistantModelMessage;
      }

      return {
        role: "user",
        content: `${username} on ${createdAt}: ${content}`,
      } as UserModelMessage;
    })
  );

  const model: LanguageModel = (() => {
    const models = [
      anthropic("claude-3-5-haiku-20241022"),
      deepSeek("deepseek-chat"),
      openaiCompatible("google/gemini-2.5-pro"),
    ];

    return models[Math.floor(Math.random() * models.length)]!;
  })();

  console.log(`using ${model.modelId}`);

  const { text } = await generateText({
    model,
    messages: [
      {
        role: "system",
        content: neroka,
      },
      {
        role: "system",
        content: `It's currently ${moment().format("MMMM Do YYYY, h:mm:ss a")}.`,
      },
      ...userMessages,
    ],
  });

  const reply = await restClient.v1.statuses.create({
    inReplyToId: status.id,
    visibility: "direct",
    status: `@${status.account.username} ${text}`,
  });

  const replyContent = textContent(reply.content);
  console.log(
    `reply is ${replyContent.substring(0, Math.min(100, replyContent.length))}`
  );
}
