package sanch.pet;

import org.telegram.telegrambots.longpolling.TelegramBotsLongPollingApplication;

public class Main {
    public static void main(String[] args) {
        try (TelegramBotsLongPollingApplication botsApplication = new TelegramBotsLongPollingApplication()) {
            botsApplication.registerBot(BotConfig.BOT_TOKEN, new ChurkaBot(BotConfig.BOT_TOKEN));
            System.out.println("ChurkaBot successfully started! Let's have some fun!!!");
            // Ensure this prcess wait forever
            Thread.currentThread().join();
        } catch (Exception e) {
            e.printStackTrace();
        }
    }
}
