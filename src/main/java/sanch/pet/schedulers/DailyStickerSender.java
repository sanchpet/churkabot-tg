package sanch.pet.schedulers;

import org.springframework.scheduling.annotation.Scheduled;
import org.springframework.stereotype.Component;
import org.telegram.telegrambots.meta.api.methods.send.SendSticker;
import org.telegram.telegrambots.meta.exceptions.TelegramApiException;
import org.telegram.telegrambots.meta.generics.TelegramClient;

import sanch.pet.services.LogStrings;
import sanch.pet.services.StickerCollection;

import lombok.extern.slf4j.Slf4j;

@Slf4j
@Component
public class DailyStickerSender {

    private final TelegramClient telegramClient;
    private final String chatId = "-4978557339"; // chat to send the sticker to

    public DailyStickerSender(TelegramClient telegramClient) {
        this.telegramClient = telegramClient;
    }

    @Scheduled(cron = "0 0 * * * ?")
    public void sendDailySticker() {
        SendSticker sendSticker = SendSticker.builder()
            .chatId(chatId)
            .sticker(StickerCollection.CHURKA_JOKER.getSticker())
            .build();
        try {
            telegramClient.execute(sendSticker);
            log.info(LogStrings.dailySticker(StickerCollection.CHURKA_JOKER.name(), chatId));
        } catch (TelegramApiException e) {
            log.error(LogStrings.dailyStickerError(StickerCollection.CHURKA_JOKER.name(), chatId), e);
        }
    }
}
