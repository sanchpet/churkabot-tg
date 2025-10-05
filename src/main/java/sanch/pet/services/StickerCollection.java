package sanch.pet.services;

import org.telegram.telegrambots.meta.api.objects.InputFile;

public enum StickerCollection {
    CHURKA_JOKER(new InputFile("CAACAgIAAxkBAAMxaNYNncvw1laDYnDFYW2N-PBvprEAAlYTAALzWwFIbKJ09BOwAAHkNgQ")),
    CHURKA_SILENCE(new InputFile("CAACAgIAAxkBAAOLaNYh6pGe8n_qvV40Hj_SP9vqQJYAAqEVAAINfRBIc9_-Nn_UjbI2BA")),
    CHURKA_PROOF(new InputFile("CAACAgIAAxkBAAIBl2jiG4Qy5rmnf9zTzIILZbanfSDOAAKXFQAC1MIRSFaAhf9LJ34gNgQ"));

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