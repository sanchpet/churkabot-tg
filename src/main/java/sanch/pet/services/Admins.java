package sanch.pet.services;

public enum Admins {
    SANCHPET("sanchpet"),
    CHURKAMAN("churkaman"),
    BOTMASTER("botmaster");

    private final String triggerWord;

    Admins(String triggerWord) {
        this.triggerWord = triggerWord;
    }

    public String getTriggerWord() {
        return triggerWord;
    }

    public static boolean contains(String username) {
        for (Admins u : values()) {
            if (u.getTriggerWord().equalsIgnoreCase(username)) {
                return true;
            }
        }
        return false;
    }
}
