package ru.bmstr.java2go;

import org.junit.jupiter.api.Test;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.context.annotation.Import;

@Import(TestcontainersConfiguration.class)
@SpringBootTest
class Java2GoApplicationTests {

    @Test
    void contextLoads() {
    }

}
