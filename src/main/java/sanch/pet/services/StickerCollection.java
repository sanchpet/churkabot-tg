package sanch.pet.services;

import org.telegram.telegrambots.meta.api.objects.InputFile;

public enum StickerCollection {
    CHURKA_JOKER(new InputFile("CAACAgIAAxkBAAMxaNYNncvw1laDYnDFYW2N-PBvprEAAlYTAALzWwFIbKJ09BOwAAHkNgQ")),
    CHURKA_SMILE(new InputFile("CAACAgIAAxkBAAIBQ2YNncvw1laDYnDFYW2N-PBvprEAAlYTAALzWwFIbKJ09BOwAAHkNgR")),
    CHURKA_SAD(new InputFile("CAACAgIAAxkBAAIBQmYNncvw1laDYnDFYW2N-PBvprEAAlYTAALzWwFIbKJ09BOwAAHkNgS"));

    private final InputFile sticker;

    StickerCollection(InputFile sticker) {
        this.sticker = sticker;
    }

    public InputFile getSticker() {
        return sticker;
    }

    @Override
    public String toString() {
        return sticker.toString();
    }
}