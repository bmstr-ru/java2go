package ru.bmstr.java2go.persistence;

import org.springframework.data.jpa.repository.Modifying;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.CrudRepository;

import java.math.BigDecimal;
import java.util.List;

public interface CurrencyRateRepository extends CrudRepository<CurrencyRateRecord, Long> {

    @Modifying
    @Query(nativeQuery = true, value = """
            insert into currency_rate (base_currency, quoted_currency, rate)
            values (:baseCurrency, :quotedCurrency, :rate)
            on conflict (base_currency, quoted_currency)
            do update set rate = :rate
            """)
    void saveRate(String baseCurrency, String quotedCurrency, BigDecimal rate);

    CurrencyRateRecord findByBaseCurrencyAndQuotedCurrency(String baseCurrency, String quotedCurrency);

    List<CurrencyRateRecord> findAll();
}
