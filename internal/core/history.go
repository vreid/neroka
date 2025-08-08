package core

/*func (n *Neroka) lookupDisplayName(userId string) string {
	user, err := n.session.User(userId)
	if err != nil || user == nil {
		return userId
	}

	displayName := user.DisplayName()
	if len(displayName) == 0 {
		return userId
	}

	return displayName
}

func (n *Neroka) AddToHistory(guildId, channelId, userId, userName, role, content string) error {
	if _, ok := n.conversations[channelId]; !ok {
		n.conversations[channelId] = true
	}

	if len(guildId) > 0 &&
		len(userId) > 0 &&
		role == "user" {
		if _, ok := n.recentParticipants[channelId]; !ok {
			n.recentParticipants[channelId] = []string{}
		}
		n.recentParticipants[channelId] = append(n.recentParticipants[channelId], userId)
	}

	user, err := n.session.User(userId)
	if err != nil {
		return err
	}

	isDm := len(guildId) == 0
	isOtherBot := false
	if user != nil &&
		user.Bot &&
		userId != n.session.State.User.ID {
		isOtherBot = true
	}

	formattedContent := ""
	if role == "user" && len(userName) > 0 {
		if isDm {
			formattedContent = content
		} else {
			cleanContent := ConvertMentionsToNames(content, n.lookupDisplayName)
			if isOtherBot {
				formattedContent = fmt.Sprintf("%s: %s", userName, cleanContent)
			} else {
				formattedContent = fmt.Sprintf("%s (<@%s>): %s", userName, userId, cleanContent)
			}
		}
	} else {
		if len(guildId) > 0 {
			formattedContent = ConvertMentionsToNames(content, n.lookupDisplayName)
		} else {
			formattedContent = content
		}
	}

	messageContent := formattedContent
	//hasImages := false

	return nil
}*/

/*
async def add_to_history(channel_id: int, role: str, content: str, user_id: int = None, guild_id: int = None, attachments: List[discord.Attachment] = None, user_name: str = None):

    # Handle image attachments - create complex content for AI providers that support images
    message_content = formatted_content
    has_images = False

    if role == "user" and attachments and not is_other_bot:
        # Get current provider to determine image format
        provider_name = "claude"  # default

        if is_dm and user_id:
            # For DMs, get provider from selected server or shared guild
            selected_guild_id = dm_server_selection.get(user_id)
            if selected_guild_id:
                provider_name, _ = ai_manager.get_guild_settings(selected_guild_id)
            else:
                # Try to get from shared guild
                shared_guild = get_shared_guild(user_id)
                if shared_guild:
                    provider_name, _ = ai_manager.get_guild_settings(shared_guild.id)
        elif guild_id:
            # For servers, get provider directly
            provider_name, _ = ai_manager.get_guild_settings(guild_id)

        # Process images if provider supports them
        if provider_name in ["claude", "gemini", "openai", "custom"]:
            image_parts = []
            text_parts = []

            # Add text content first
            if formatted_content.strip():
                if provider_name == "openai":
                    text_parts.append({"type": "text", "text": formatted_content})
                else:
                    text_parts.append({"type": "text", "text": formatted_content})

            # Process each attachment
            for attachment in attachments:
                if any(attachment.filename.lower().endswith(ext) for ext in ['.png', '.jpg', '.jpeg', '.gif', '.webp']):
                    # Size check
                    if attachment.size > 20 * 1024 * 1024:  # 20MB limit for OpenAI, 3MB for others
                        size_limit = "20MB" if provider_name == "openai" else "30MB"
                        text_parts.append({"type": "text", "text": f" [Image {attachment.filename} was too large (limit: {size_limit})]"})
                        continue

                    try:
                        image_data = await process_image_attachment(attachment, provider_name)
                        if image_data:
                            image_parts.append(image_data)
                            has_images = True
                        else:
                            text_parts.append({"type": "text", "text": f" [Could not process image {attachment.filename}]"})
                    except Exception as e:
                        print(f"Error processing image {attachment.filename}: {e}")
                        text_parts.append({"type": "text", "text": f" [Error processing image {attachment.filename}]"})
                else:
                    # Non-image attachment
                    text_parts.append({"type": "text", "text": f" [File: {attachment.filename}]"})

            # Combine text and images into complex content
            if has_images:
                message_content = text_parts + image_parts
            else:
                # No valid images, use regular text content with attachment notes
                attachment_notes = []
                for attachment in attachments:
                    attachment_notes.append(f"[Attachment: {attachment.filename}]")

                if attachment_notes:
                    message_content = formatted_content + " " + " ".join(attachment_notes)
        else:
            # Provider doesn't support images, add text descriptions
            attachment_parts = []
            for attachment in attachments:
                if any(attachment.filename.lower().endswith(ext) for ext in ['.png', '.jpg', '.jpeg', '.gif', '.webp']):
                    attachment_parts.append(f"[Image: {attachment.filename} - current AI model doesn't support images]")
                else:
                    attachment_parts.append(f"[File: {attachment.filename}]")

            if attachment_parts:
                message_content = formatted_content + " " + " ".join(attachment_parts)

    # Check if we should group with the previous message (only for text content)
    should_group = False
    if conversations[channel_id] and not has_images:  # Don't group messages with images
        last_message = conversations[channel_id][-1]

        if (last_message["role"] == role and
            isinstance(last_message["content"], str) and  # Only group with text messages
            ((role == "user" and user_id and
            last_message["content"].startswith(f"{user_name} (<@{user_id}>):")) or
            (role == "assistant"))):
            should_group = True

    if should_group and role == "user" and isinstance(message_content, str):
        # Group with previous user message
        if isinstance(conversations[channel_id][-1]["content"], str):
            existing_content = conversations[channel_id][-1]["content"] or ""
            new_content = content or ""
            conversations[channel_id][-1]["content"] = existing_content + f"\n{new_content}"
        else:
            # Don't group if previous message has complex content
            conversations[channel_id].append({"role": role, "content": message_content})
    else:
        # Create new message entry
        conversations[channel_id].append({"role": role, "content": message_content})

    # Maintain history length limit
    if is_dm:
        selected_guild_id = dm_server_selection.get(user_id) if user_id else None
        if selected_guild_id:
            max_history = get_history_length(selected_guild_id)
        elif guild_id:
            max_history = get_history_length(guild_id)
        else:
            max_history = 50
    else:
        max_history = get_history_length(guild_id) if guild_id else 50

    if len(conversations[channel_id]) > max_history:
        conversations[channel_id] = conversations[channel_id][-max_history:]
*/
