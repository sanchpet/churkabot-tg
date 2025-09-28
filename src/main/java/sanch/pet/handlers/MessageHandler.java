package sanch.pet.handlers;

import java.io.InvalidObjectException;
import java.util.Comparator;
import java.util.List;

import org.telegram.telegrambots.meta.api.methods.send.SendMessage;
import org.telegram.telegrambots.meta.api.methods.send.SendPhoto;
import org.telegram.telegrambots.meta.api.methods.send.SendSticker;
import org.telegram.telegrambots.meta.api.objects.InputFile;
import org.telegram.telegrambots.meta.api.objects.message.Message;
import org.telegram.telegrambots.meta.api.objects.photo.PhotoSize;
import org.telegram.telegrambots.meta.api.objects.replykeyboard.InlineKeyboardMarkup;
import org.telegram.telegrambots.meta.api.objects.replykeyboard.buttons.InlineKeyboardButton;
import org.telegram.telegrambots.meta.api.objects.replykeyboard.buttons.InlineKeyboardRow;
import org.telegram.telegrambots.meta.api.objects.stickers.Sticker;
import org.telegram.telegrambots.meta.exceptions.TelegramApiException;
import org.telegram.telegrambots.meta.generics.TelegramClient;

import sanch.pet.services.Emoji;
import sanch.pet.services.StickerCollection;
import sanch.pet.services.TriggerWords;
import sanch.pet.services.LogStrings;
import sanch.pet.services.BotUtils;

import lombok.extern.slf4j.Slf4j;

@Slf4j
public class MessageHandler {
    public static void HandleIncomingMessage(Message message, TelegramClient telegramClient) throws InvalidObjectException {
        if (message == null) return;

        boolean isGroup = message.getChat().isGroupChat() || message.getChat().isSuperGroupChat();

        if (isGroup) {
            handleGroupChat(message, telegramClient);
        } else {
            handlePrivateChat(message, telegramClient);
        }
    }

    public static void handleGroupChat(Message message, TelegramClient telegramClient) throws InvalidObjectException {
        String user_first_name = message.getFrom().getFirstName();
        String user_last_name = message.getFrom().getLastName();
        String user_username = message.getFrom().getUserName();
        Integer message_id = message.getMessageId();
        String user_id = String.valueOf(message.getFrom().getId());
        String message_text = message.getText();
        String chat_id = String.valueOf(message.getChatId());

        boolean hasTrigger = TriggerWords.containsTriggerWord(message_text);
        if (!hasTrigger) return;

        SendSticker sendSticker = SendSticker.builder()
            .chatId(chat_id)
            .sticker(StickerCollection.CHURKA_JOKER.getSticker())
            .replyToMessageId(message_id)
            .build();
        try {
            telegramClient.execute(sendSticker);
            log.info(LogStrings.stickerReaction(StickerCollection.CHURKA_JOKER.name(), user_first_name, user_last_name, user_id, user_username, chat_id));
        } catch (TelegramApiException e) {
            log.error(LogStrings.stickerReactionError(StickerCollection.CHURKA_JOKER.name(), user_first_name, user_last_name, user_id, user_username, chat_id), e);
        }
    }

    public static void handlePrivateChat(Message message, TelegramClient telegramClient) throws InvalidObjectException {
        String user_first_name = message.getFrom().getFirstName();
        String user_last_name = message.getFrom().getLastName();
        String user_username = message.getFrom().getUserName();
        Integer message_id = message.getMessageId();
        String user_id = String.valueOf(message.getFrom().getId());
        String message_text = message.getText();
        String reply_text = message_text + " " + Emoji.GRINNING_FACE_WITH_SMILING_EYES;
        String chat_id = String.valueOf(message.getChatId());

        if (message.hasText()) {
            boolean isAdmin = BotUtils.isPrivateChatFromAdmin(message);
            if (!isAdmin) {
                reply_text = "Sorry, chat only for admins.";
                SendMessage answer = SendMessage // Create a message object
                    .builder()
                    .chatId(chat_id)
                    .text(reply_text)
                    .build();
                try {
                    telegramClient.execute(answer);
                    log.info("Non-admin user " + user_username + " (" + user_first_name + " " + user_last_name + " - id = " + user_id + ") tried to access the bot in private chat.");
                } catch (TelegramApiException e) {
                    log.error("Error sending message to non-admin user " + user_username + " (" + user_first_name + " " + user_last_name + " - id = " + user_id + ")", e);
                }
                return;
            }

            if (message_text != null && message_text.equals("/start")) {
                reply_text = "Привет, " + user_first_name + "! Я чуркабот";
            } else if (message_text != null && message_text.equals("/help")) {
                reply_text = "На текущем этапе тебе ничего не поможет. Обратись к @sanchpet_unfiltered.";
            } else {
                boolean hasTrigger = TriggerWords.containsTriggerWord(message_text);
                if (hasTrigger) {
                    SendSticker sendSticker = SendSticker.builder()
                        .chatId(chat_id)
                        .sticker(StickerCollection.CHURKA_JOKER.getSticker())
                        .replyToMessageId(message_id)
                        .build();
                    try {
                        telegramClient.execute(sendSticker);
                        log.info(LogStrings.stickerReaction(StickerCollection.CHURKA_JOKER.name(), user_first_name, user_last_name, user_id, user_username, chat_id));
                    } catch (TelegramApiException e) {
                        log.info(LogStrings.stickerReactionError(StickerCollection.CHURKA_JOKER.name(), user_first_name, user_last_name, user_id, user_username, chat_id));
                    }
                } else {
                    SendMessage answer = SendMessage // Create a message object
                        .builder()
                        .chatId(chat_id)
                        .text(reply_text)
                        .replyMarkup(InlineKeyboardMarkup
                                    .builder()
                                    .keyboardRow(
                                            new InlineKeyboardRow(InlineKeyboardButton
                                                    .builder()
                                                    .text("Update message text")
                                                    .callbackData("update_msg_text")
                                                    .build()
                                            )
                                    )
                                    .build())
                        .build();
                    try {
                        telegramClient.execute(answer);
                        log.info("Replied to admin " + user_username + " (" + user_first_name + " " + user_last_name + " - id = " + user_id + ") in private chat.");
                    } catch (TelegramApiException e) {
                        log.error("Error sending message to admin " + user_username + " (" + user_first_name + " " + user_last_name + " - id = " + user_id + ") in private chat.", e);
                    }
                }
            }
        } else if (message.hasPhoto()) {
            List<PhotoSize> photos = message.getPhoto();
            // Know file_id
            String f_id = photos.stream().max(Comparator.comparing(PhotoSize::getFileSize))
                .map(PhotoSize::getFileId)
                .orElse("");
            // Know photo width
            int f_width = photos.stream().max(Comparator.comparing(PhotoSize::getFileSize))
                .map(PhotoSize::getWidth)
                .orElse(0);
            // Know photo height
            int f_height = photos.stream().max(Comparator.comparing(PhotoSize::getFileSize))
                .map(PhotoSize::getHeight)
                .orElse(0);
            // Set photo caption
            String caption = "file_id: " + f_id + "\nwidth: " + Integer.toString(f_width) + "\nheight: " + Integer.toString(f_height);

            SendPhoto answer = SendPhoto // Create a message object
                .builder()
                .chatId(chat_id)
                .photo(new InputFile(f_id))
                .caption(caption)
                .build();
            try {
                telegramClient.execute(answer);
                log.info("User sent a photo with file ID: " + f_id);
            } catch (TelegramApiException e) {
                log.error("Error sending photo back to user", e);
            }
        } else if (message.hasSticker()) {
            Sticker sticker = message.getSticker();
            String stickerId = sticker.getFileId();
            SendMessage answer = SendMessage
                .builder()
                .chatId(chat_id)
                .text("Вы мне отправили стикер с ID: " + stickerId)
                .build();
            try {
                telegramClient.execute(answer); // Sending our message object to user
                log.info("User sent a sticker with file ID: " + stickerId);
            } catch (TelegramApiException e) {
                log.error("Error sending sticker ID back to user", e);
            }
        }
    }
}
