import {
  generateText,
  type AssistantModelMessage,
  type LanguageModel,
  type UserModelMessage,
} from "ai";
import { convert } from "html-to-text";
import type { MediaAttachment } from "masto/mastodon/entities/v1/media-attachment.js";
import type { Status } from "masto/mastodon/entities/v1/status.js";
import moment from "moment";
import { EOL } from "node:os";
import sharp from "sharp";

export function textContent(content: string): string {
  return convert(content, {
    wordwrap: false,
    selectors: [{ selector: "a", format: "skip" }],
  });
}

export async function describeAttachment(
  visionModel: LanguageModel,
  mediaAttachment: MediaAttachment
) {
  if (mediaAttachment.type != "image") {
    console.log(`media attachment ${mediaAttachment.id} wasn't an image`);
    return "";
  }

  if (mediaAttachment.url == null) {
    console.log(`media attachment ${mediaAttachment.id} had no url`);
    return "";
  }

  console.log(`describing media attachment ${mediaAttachment.id}`);

  const response = await fetch(mediaAttachment.url);
  const buffer = await sharp(await response.arrayBuffer())
    .resize(1000)
    .webp()
    .toBuffer();

  const { text } = await generateText({
    model: visionModel,
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
    `description of media attachment is ${text.substring(0, Math.min(100, text.length))}`
  );
  return `[${text}]`;
}

export async function generateImageDescriptions(
  visionModel: LanguageModel,
  status: Status
): Promise<string> {
  return (
    await Promise.all(
      status.mediaAttachments.map(
        async (mediaAttachment) =>
          await describeAttachment(visionModel, mediaAttachment)
      )
    )
  )
    .filter((x) => x && x.length > 0)
    .join(EOL);
}

export async function convertToMessages(
  visionModel: LanguageModel,
  conversation: Status[],
  botUsername: string
) {
  return await Promise.all(
    conversation.map(async (entry) => {
      const username = entry.account.username;

      const createdAt = moment(entry.createdAt).format(
        "MMMM Do YYYY, h:mm:ss a"
      );

      const imageDescriptions = await generateImageDescriptions(
        visionModel,
        entry
      );

      const content = (() => {
        if (imageDescriptions.length == 0) {
          return `${textContent(entry.content)}`;
        }

        return [`${textContent(entry.content)}`, imageDescriptions].join(EOL);
      })();

      if (botUsername && username == botUsername) {
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
}

export async function convertToDialog(
  visionModel: LanguageModel,
  conversation: Status[]
) {
  return await Promise.all(
    conversation.map(async (entry) => {
      const username = entry.account.username;

      const createdAt = moment(entry.createdAt).format(
        "MMMM Do YYYY, h:mm:ss a"
      );

      const imageDescriptions = await generateImageDescriptions(
        visionModel,
        entry
      );

      const content = (() => {
        if (imageDescriptions.length == 0) {
          return `${textContent(entry.content)}`;
        }

        return [`${textContent(entry.content)}`, imageDescriptions].join(EOL);
      })();

      return `${username} on ${createdAt}: ${content}`;
    })
  );
}
