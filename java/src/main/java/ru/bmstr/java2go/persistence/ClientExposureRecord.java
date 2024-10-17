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
@Table(name = "client_exposure")
@Getter
@Setter
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class ClientExposureRecord {
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;
    private Long clientId;
    private BigDecimal totalExposureAmount;
    @Column(length = 3)
    @JdbcType(CharJdbcType.class)
    private String totalExposureCurrency;

}
