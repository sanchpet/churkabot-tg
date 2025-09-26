package sanch.pet.handlers;

import java.io.InvalidObjectException;
import java.text.DateFormat;
import java.text.SimpleDateFormat;
import java.util.Comparator;
import java.util.Date;
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
import lombok.extern.slf4j.Slf4j;

@Slf4j
public class MessageHandler {
    public static void HandleIncomingMessage(Message message, TelegramClient telegramClient) throws InvalidObjectException {
        if (message == null) return;

        String user_first_name = message.getChat().getFirstName();
        String user_last_name = message.getChat().getLastName();
        String user_username = message.getChat().getUserName();
        long user_id = message.getChat().getId();
        String message_text = message.getText();
        String reply_text = message_text + " " + Emoji.GRINNING_FACE_WITH_SMILING_EYES;
        long chat_id = message.getChatId();

        boolean isGroup = message.getChat().isGroupChat() || message.getChat().isSuperGroupChat();

        if (isGroup) {
            boolean hasTrigger = false;
            if (message_text != null) {
                for (TriggerWords tw : TriggerWords.values()) {
                    if (message_text.toLowerCase().contains(tw.toString().toLowerCase())) {
                        hasTrigger = true;
                        break;
                    }
                }
            }
            if (!hasTrigger) return;
        }

        if (message.hasText()) {
            if (message_text != null && message_text.equals("/start")) {
                reply_text = "Hi, " + user_first_name + "! This is ChurkaBot! Send me a photo and I will reply you with its file_id, width and height!";
            } else if (message_text != null && message_text.equals("/help")) {
                reply_text = "This bot was created to help you get file_id, width and height of photos. Just send me a photo and I will reply you with this information.";
            } else {
                boolean hasTrigger = TriggerWords.containsTriggerWord(message_text);
                if (hasTrigger) {
                    SendSticker sendSticker = SendSticker.builder()
                        .chatId(chat_id)
                        .sticker(StickerCollection.CHURKA_JOKER.getSticker())
                        .build();
                    try {
                        telegramClient.execute(sendSticker);
                    } catch (TelegramApiException e) {
                        e.printStackTrace();
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
                        telegramClient.execute(answer); // Sending our message object to user
                        log.info(logstring(user_first_name, user_last_name, Long.toString(user_id), user_username, message_text, reply_text));
                    } catch (TelegramApiException e) {
                        e.printStackTrace();
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
                telegramClient.execute(answer); // Sending our message object to user
                log.info(logstring(user_first_name, user_last_name, Long.toString(user_id), user_username, message_text, caption));
            } catch (TelegramApiException e) {
                e.printStackTrace();
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
                e.printStackTrace();
            }
        }
    }

    private static String logstring(String first_name, String last_name, String user_id, String user_name, String txt, String bot_answer) {
        DateFormat dateFormat = new SimpleDateFormat("yyyy/MM/dd HH:mm:ss");
        Date date = new Date();

        return "\n----------------------------\n" +
            dateFormat.format(date) + "\n" +
            String.format("Message from %s (%s %s). (id = %s) \nText - %s\n", user_name, first_name, last_name, user_id, txt) +
            String.format("Bot answer: \nText - %s\n", bot_answer);
    }
}
