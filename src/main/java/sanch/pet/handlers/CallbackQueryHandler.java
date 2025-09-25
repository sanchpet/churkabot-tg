package sanch.pet.handlers;

import java.io.InvalidObjectException;

import org.telegram.telegrambots.meta.api.methods.updatingmessages.EditMessageText;
import org.telegram.telegrambots.meta.api.objects.CallbackQuery;
import org.telegram.telegrambots.meta.exceptions.TelegramApiException;
import org.telegram.telegrambots.meta.generics.TelegramClient;

import sanch.pet.services.Emoji;

public class CallbackQueryHandler {
    public static void HandleIncomingCallbackQuery(CallbackQuery callbackQuery, TelegramClient telegramClient) throws InvalidObjectException {
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
}
