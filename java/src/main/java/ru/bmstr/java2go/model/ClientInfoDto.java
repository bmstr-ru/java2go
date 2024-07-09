package ru.bmstr.java2go.model;

public record ClientInfoDto(
        Long id,
        String name,
        Address address,
        String baseCurrency
) {
    public record Address(
            String country,
            String region,
            String zipCode,
            String address
    ){}
}
