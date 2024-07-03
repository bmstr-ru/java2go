package ru.bmstr.java2go;

import org.springframework.boot.SpringApplication;

public class TestJava2GoApplication {

    public static void main(String[] args) {
        SpringApplication.from(Java2GoApplication::main).with(TestcontainersConfiguration.class).run(args);
    }

}
