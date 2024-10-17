package ru.bmstr.java2go.persistence;

import jakarta.persistence.Column;
import jakarta.persistence.Entity;
import jakarta.persistence.GeneratedValue;
import jakarta.persistence.GenerationType;
import jakarta.persistence.Id;
import jakarta.persistence.Table;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;
import org.hibernate.annotations.JdbcType;
import org.hibernate.type.descriptor.jdbc.CharJdbcType;

import java.math.BigDecimal;

@Entity
@Table(name = "deal")
@Getter
@Setter
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class DealRecord {
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;
    private Long dealId;
    private Long clientId;
    private BigDecimal boughtAmount;
    @Column(length = 3)
    @JdbcType(CharJdbcType.class)
    private String boughtCurrency;
    private BigDecimal soldAmount;
    @Column(length = 3)
    @JdbcType(CharJdbcType.class)
    private String soldCurrency;
}
