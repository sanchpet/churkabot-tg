package sanch.pet.updateshandlers;

import org.telegram.telegrambots.client.okhttp.OkHttpTelegramClient;
import org.telegram.telegrambots.longpolling.util.LongPollingSingleThreadUpdateConsumer;
import org.telegram.telegrambots.meta.api.methods.send.SendMessage;
import org.telegram.telegrambots.meta.api.methods.send.SendPhoto;
import org.telegram.telegrambots.meta.api.objects.CallbackQuery;
import org.telegram.telegrambots.meta.api.objects.InputFile;
import org.telegram.telegrambots.meta.api.objects.Update;
import org.telegram.telegrambots.meta.exceptions.TelegramApiException;
import org.telegram.telegrambots.meta.generics.TelegramClient;
import org.telegram.telegrambots.meta.api.objects.message.Message;
import org.telegram.telegrambots.meta.api.objects.photo.PhotoSize;
import org.telegram.telegrambots.meta.api.objects.replykeyboard.InlineKeyboardMarkup;
import org.telegram.telegrambots.meta.api.objects.replykeyboard.buttons.InlineKeyboardButton;
import org.telegram.telegrambots.meta.api.objects.replykeyboard.buttons.InlineKeyboardRow;
import org.telegram.telegrambots.meta.api.methods.updatingmessages.EditMessageText;

import sanch.pet.services.Emoji;

import java.io.InvalidObjectException;
import java.text.DateFormat;
import java.text.SimpleDateFormat;
import java.util.Comparator;
import java.util.Date;
import java.util.List;

public class ChurkaBot implements LongPollingSingleThreadUpdateConsumer {
    private final TelegramClient telegramClient;

    public ChurkaBot(String botToken) {
        telegramClient = new OkHttpTelegramClient(botToken);
    }

    @Override
    public void consume(Update update) {
        if (update.hasMessage()) {
            try {
                Message message = update.getMessage();
                try {
                    handleIncomingMessage(message);
                } catch (InvalidObjectException e) {
                    e.printStackTrace();
                }
            } catch (Exception e) {
                e.printStackTrace();
            }
        } else if (update.hasCallbackQuery()) {
            try {
                CallbackQuery callbackQuery = update.getCallbackQuery();
                try {
                    handleIncomingCallbackQuery(callbackQuery);
                } catch (InvalidObjectException e) {
                    e.printStackTrace();
                }
            } catch (Exception e) {
                e.printStackTrace();
            }
            
        }
    }   

    private void handleIncomingMessage(Message message) throws InvalidObjectException {
        if (message == null) return;

        String user_first_name = message.getChat().getFirstName();
        String user_last_name = message.getChat().getLastName();
        String user_username = message.getChat().getUserName();
        long user_id = message.getChat().getId();
        String message_text = message.getText();
        String reply_text = message_text + " " + Emoji.GRINNING_FACE_WITH_SMILING_EYES;
        long chat_id = message.getChatId();

        if (message.hasText()) {
            if (message_text.equals("/start")) {
                reply_text = "Hi, " + user_first_name + "! This is ChurkaBot! Send me a photo and I will reply you with its file_id, width and height!";
            } else if (message_text.equals("/help")) {
                reply_text = "This bot was created to help you get file_id, width and height of photos. Just send me a photo and I will reply you with this information.";
            } else {
                reply_text = "I don't understand you. Please, send me a photo so I can get its file_id, width and height.";
            }
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
                log(user_first_name, user_last_name, Long.toString(user_id), user_username, message_text, reply_text);
            } catch (TelegramApiException e) {
                e.printStackTrace();
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
                log(user_first_name, user_last_name, Long.toString(user_id), user_username, message_text, caption);
            } catch (TelegramApiException e) {
                e.printStackTrace();
            }
        }
    }

    private void handleIncomingCallbackQuery(CallbackQuery callbackQuery) throws InvalidObjectException {
        // This method is not used in this example
        String callbackData = callbackQuery.getData();
        long chat_id = callbackQuery.getMessage().getChatId();
        int message_id = callbackQuery.getMessage().getMessageId();
        if (callbackData.equals("update_msg_text")) {
            String newText = "Message text updated " + Emoji.GRINNING_FACE_WITH_SMILING_EYES;
            EditMessageText new_message = EditMessageText
                .builder()
                .chatId(chat_id)
                .messageId(message_id)
                .text(newText)
                .build();
            try {
                telegramClient.execute(new_message);
            } catch (TelegramApiException e) {
                e.printStackTrace();
            }
        }
    }

    private void log(String first_name, String last_name, String user_id, String user_name, String txt, String bot_answer) {
        System.out.println("\n ----------------------------");
        DateFormat dateFormat = new SimpleDateFormat("yyyy/MM/dd HH:mm:ss");
        Date date = new Date();
        System.out.println(dateFormat.format(date));
        System.out.println("Message from " + user_name + " (" + first_name + " " + last_name + ")" + ". (id = " + user_id + ") \n Text - " + txt);
        System.out.println("Bot answer: \n Text - " + bot_answer);
    }
}
