package sanch.pet.services;

import org.telegram.telegrambots.meta.api.objects.message.Message;
import org.telegram.telegrambots.meta.api.objects.chat.Chat;

public class BotUtils {
    public static boolean isPrivateChatFromAdmin(Message message) {
        Chat chat = message.getChat();
        if (chat.isUserChat()) {  // Checks if chat type is private (user chat)
            String username = message.getFrom().getUserName();
            return Admins.contains(username);
        }
        return false;
    }
}
