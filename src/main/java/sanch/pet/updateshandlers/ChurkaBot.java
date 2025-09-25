package sanch.pet.updateshandlers;

import org.telegram.telegrambots.client.okhttp.OkHttpTelegramClient;
import org.telegram.telegrambots.longpolling.util.LongPollingSingleThreadUpdateConsumer;
import org.telegram.telegrambots.meta.api.objects.CallbackQuery;
import org.telegram.telegrambots.meta.api.objects.Update;
import org.telegram.telegrambots.meta.generics.TelegramClient;
import org.telegram.telegrambots.meta.api.objects.message.Message;
import sanch.pet.handlers.MessageHandler;
import sanch.pet.handlers.CallbackQueryHandler;

import java.io.InvalidObjectException;

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
                    MessageHandler.HandleIncomingMessage(message, telegramClient);
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
                    CallbackQueryHandler.HandleIncomingCallbackQuery(callbackQuery, telegramClient);
                } catch (InvalidObjectException e) {
                    e.printStackTrace();
                }
            } catch (Exception e) {
                e.printStackTrace();
            }
            
        }
    }   
}
