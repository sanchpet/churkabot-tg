package sanch.pet.services;

public class LogStrings {
    public static String stickerReaction(String sticker, String first_name, String last_name, String user_id, String user_name, String chat_id) {
        return String.format("Sent sticker %s as reaction to message from %s (%s %s - id = %s) in chat %s", 
                            sticker, user_name, first_name, last_name, user_id, chat_id);
    }

    public static String stickerReactionError(String sticker, String first_name, String last_name, String user_id, String user_name, String chat_id) {
        return String.format("Error while sending sticker %s as reaction to message from %s (%s %s - id = %s) in chat %s:", 
                            sticker, user_name, first_name, last_name, user_id, chat_id);
    }

    public static String dailySticker(String sticker, String chat_id) {
        return String.format("Sent daily sticker %s in chat %s", 
                            sticker, chat_id);
    }

    public static String dailyStickerError(String sticker, String chat_id) {
        return String.format("Error while sending daily sticker %s in chat %s:", 
                            sticker, chat_id);
    }
}
