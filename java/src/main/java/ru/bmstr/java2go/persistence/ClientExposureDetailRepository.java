package ru.bmstr.java2go.persistence;

import org.springframework.data.repository.CrudRepository;

import java.util.List;
import java.util.Optional;

public interface ClientExposureDetailRepository extends CrudRepository<ClientExposureDetailRecord, Long> {

    Optional<ClientExposureDetailRecord> findByClientIdAndExposureCurrency(Long clientId, String exposureCurrency);

    List<ClientExposureDetailRecord> findAllByClientId(Long clientId);
}
