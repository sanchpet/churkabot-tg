package sanch.pet.services;

public enum Admins {
    SANCHPET("sanchpet"),
    AMPLITUDNIY("amplitudniy"),
    SANCHPET_UNFILTERED("sanchpet_unfiltered");

    private final String username;

    Admins(String username) {
        this.username = username;
    }

    public String getUsername() {
        return username;
    }

    public static boolean contains(String username) {
        if (username == null) return false;
        for (Admins u : values()) {
            if (u.getUsername().equalsIgnoreCase(username)) {
                return true;
            }
        }
        return false;
    }
}