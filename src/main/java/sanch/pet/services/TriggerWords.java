package sanch.pet.services;

import java.util.regex.Pattern;
import java.util.regex.Matcher;

public enum TriggerWords {
    CHUR("чур(ка|ки|ке|кам|ками|ках|ок)?", StickerCollection.CHURKA_JOKER),
    CHUROCHK("чурочк(а|и|е|ам|ами|ах|ек)?", StickerCollection.CHURKA_JOKER),
    SILENCE("ти(ш(е|ина|иной|шком|шами|шах)?|х(о)?)", StickerCollection.CHURKA_SILENCE),
    SILENCE2("молч(ать|ит|ишь|а|и)?", StickerCollection.CHURKA_SILENCE),
    SILENCE3("заткн(ись|итесь|ешься|ёшься)?", StickerCollection.CHURKA_SILENCE),
    SILENCE4("угомон(ись|итесь|иться)?", StickerCollection.CHURKA_SILENCE),
    SILENCE5("шум(еть|ит|ишь|а|и)?", StickerCollection.CHURKA_SILENCE),
    PROOF("доказательств(о|ом|ами|а)?", StickerCollection.CHURKA_PROOF),
    PROOF2("подтвержден(о|а|ы)?", StickerCollection.CHURKA_PROOF),
    PROOF3("факт(ы|а|ом|ам|ах|ами)?", StickerCollection.CHURKA_PROOF),
    PROOF4("уверен(ы|а)?", StickerCollection.CHURKA_PROOF),
    PROOF5("докаж(и|ите|у|ет|ешь)?", StickerCollection.CHURKA_PROOF),
    PROOF7("пруф(ы|а|ом|ам|ах|ами|ай|ани)?", StickerCollection.CHURKA_PROOF),
    PROOF8("вер(ю|ишь|им|но|ят)+", StickerCollection.CHURKA_PROOF),
    PROOF6("точн(о|ый|як)+", StickerCollection.CHURKA_PROOF);

    private final Pattern pattern;
    private final StickerCollection sticker;

    TriggerWords(String regex, StickerCollection sticker) {
        this.pattern = Pattern.compile(regex, Pattern.CASE_INSENSITIVE | Pattern.UNICODE_CASE);
        this.sticker = sticker;
    }

    public boolean matches(String message) {
        if (message == null) return false;
        Matcher matcher = pattern.matcher(message);
        return matcher.find();
    }

    public StickerCollection getSticker() {
        return sticker;
    }

    public static StickerCollection findStickerFor(String message) {
        if (message == null) return null;
        for (TriggerWords trigger : values()) {
            if (trigger.matches(message)) {
                return trigger.getSticker();
            }
        }
        if (hasFiveConsecutiveUppercaseLetters(message)) {
            return StickerCollection.CHURKA_SILENCE;
        }
        return null;
    }

    public static boolean hasFiveConsecutiveUppercaseLetters(String input) {
        if (input == null || input.isEmpty()) {
            return false;
        }

        // Remove all characters except English and Russian letters
        String cleaned = input.replaceAll("[^A-Za-z\u0400-\u04FF]", "");

        int consecutiveCount = 0;
        for (int i = 0; i < cleaned.length(); i++) {
            char c = cleaned.charAt(i);
            if (Character.isUpperCase(c)) {
                consecutiveCount++;
                if (consecutiveCount >= 5) {
                    return true;
                }
            } else {
                consecutiveCount = 0;
            }
        }

        return false;
    }
}
