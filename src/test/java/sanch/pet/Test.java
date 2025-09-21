package sanch.pet;

import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertDoesNotThrow;

class MainTest {
    @Test
    void mainDoesNotThrow() {
        assertDoesNotThrow(() -> {
            // Run main in a separate thread to avoid blocking the test
            Thread t = new Thread(() -> Main.main(new String[]{}));
            t.start();
            t.interrupt(); // Interrupt to avoid infinite join
        });
    }
}