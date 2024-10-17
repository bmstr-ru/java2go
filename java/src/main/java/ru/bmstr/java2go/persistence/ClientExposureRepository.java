package ru.bmstr.java2go.persistence;

import org.springframework.data.jpa.repository.Modifying;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.CrudRepository;

import java.math.BigDecimal;
import java.util.Optional;

public interface ClientExposureRepository extends CrudRepository<ClientExposureRecord, Long> {

    @Modifying
    @Query(nativeQuery = true, value = """
            insert into client_exposure (client_id, total_exposure_amount, total_exposure_currency)
            values (:clientId, :totalExposureAmount, :totalExposureCurrency)
            on conflict (client_id)
            do update set total_exposure_amount = :totalExposureAmount, total_exposure_currency = :totalExposureCurrency
            """)
    void saveClientExposure(Long clientId, BigDecimal totalExposureAmount, String totalExposureCurrency);

    Optional<ClientExposureRecord> findByClientId(Long clientId);
}
