package sanch.pet;

import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.telegram.telegrambots.client.okhttp.OkHttpTelegramClient;
import org.telegram.telegrambots.meta.generics.TelegramClient;

@Configuration
public class BotConfig {
    @Bean
    public TelegramClient telegramClient() {
        return new OkHttpTelegramClient(System.getenv("CHURKA_BOT_TOKEN"));
    }
}
