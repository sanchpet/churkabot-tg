package sanch.pet.services;

import java.util.regex.Pattern;
import java.util.regex.Matcher;

public enum TriggerWords {
    CHUR("чур(ка|ки|ке|кам|ками|ках|ок)?"),
    CHUROCHK("чурочк(а|и|е|ам|ами|ах|ек)?");

    private final Pattern pattern;

    TriggerWords(String regex) {
        this.pattern = Pattern.compile(regex, Pattern.CASE_INSENSITIVE | Pattern.UNICODE_CASE);
    }

    public boolean matches(String message) {
        if (message == null) return false;
        Matcher matcher = pattern.matcher(message);
        return matcher.find();
    }

    public static boolean containsTriggerWord(String message) {
        if (message == null) return false;
        for (TriggerWords trigger : values()) {
            if (trigger.matches(message)) {
                return true;
            }
        }
        return false;
    }
}
