package sanch.pet.services;

public enum TriggerWords {
    CHURKA("чурка"),
    CHURKI("чурки"),
    CHURKE("чурке"),
    CHURKAM("чуркам"),
    CHURKAMI("чурками"),
    CHURKAH("чурках"),
    CHURKOV("чурок"),
    CHURKAMY("чурками"),
    CHURKACH("чурках");

    private final String username;

    TriggerWords(String username) {
        this.username = username;
    }

    @Override
    public String toString() {
        return username;
    }

    // Static method to check if the input message contains any trigger word
    public static boolean containsTriggerWord(String message) {
        if (message == null) {
            return false;
        }
        String lowerMessage = message.toLowerCase();

        for (TriggerWords trigger : TriggerWords.values()) {
            if (lowerMessage.contains(trigger.toString().toLowerCase())) {
                return true;
            }
        }
        return false;
    }
}
