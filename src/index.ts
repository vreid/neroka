import { createAnthropic } from "@ai-sdk/anthropic";
import { createDeepSeek } from "@ai-sdk/deepseek";
import { createMistral } from "@ai-sdk/mistral";
import { createOpenAICompatible } from "@ai-sdk/openai-compatible";
import { generateText, type LanguageModel } from "ai";
import { createRestAPIClient, createStreamingAPIClient } from "masto";
import moment from "moment";
import yargs from "yargs";
import { convertToMessages, textContent } from "./util";

import neroka from "./neroka.txt" with { type: "text" };
import systemPrompt from "./system_prompt.txt" with { type: "text" };

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

const openrouter = createOpenAICompatible({
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

console.log("waiting for direct messages");
for await (const event of streamingClient.public.subscribe()) {
  console.log(`received ${event.event}`);
  if (event.event != "conversation") {
    continue;
  }

  const lastStatus = event.payload.lastStatus;
  if (!lastStatus) {
    console.log("last status is null");
    continue;
  }

  const statusId = lastStatus.id;
  console.log(`status id is ${statusId}`);

  const status = await restClient.v1.statuses.$select(statusId).fetch();
  console.log("fetched status");

  if (status.account.username == credentials.username) {
    console.log(`last post is from ${credentials.username}`);
    continue;
  }

  const statusContext = await restClient.v1.statuses
    .$select(statusId)
    .context.fetch();
  console.log(
    `fetched context (${statusContext.ancestors.length} ancestors, ${statusContext.descendants.length} descendants)`
  );

  const currentConversation = [
    ...statusContext.ancestors,
    status,
    ...statusContext.descendants,
  ];

  const userMessages = await convertToMessages(
    mistral("pixtral-12b-latest"),
    currentConversation,
    credentials.username
  );

  const model: LanguageModel = (() => {
    const models = [
      anthropic("claude-3-5-haiku-20241022"),
      deepSeek("deepseek-chat"),
      mistral("mistral-medium-latest"),
      openrouter("google/gemini-2.5-pro"),
      openrouter("moonshotai/kimi-k2"),
    ];

    return models[Math.floor(Math.random() * models.length)]!;
  })();

  const { text } = await generateText({
    model,
    messages: [
      {
        role: "system",
        content: systemPrompt,
      },
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
  console.log(`finished text generation using ${model.modelId}`);

  const reply = await restClient.v1.statuses.create({
    inReplyToId: status.id,
    visibility: "direct",
    status: `@${status.account.username} ${text}`,
  });

  const replyContent = textContent(reply.content);
  console.log(
    `reply content is ${replyContent.substring(0, Math.min(100, replyContent.length))}`
  );
}
