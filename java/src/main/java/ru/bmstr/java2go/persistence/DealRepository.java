package ru.bmstr.java2go.persistence;

import org.springframework.data.repository.CrudRepository;

import java.util.List;

public interface DealRepository extends CrudRepository<DealRecord, Long> {

    List<DealRecord> findAllByClientId(Long clientId);
}
