package com.guardian.gaia.repository;

import com.guardian.custodian.Entity;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import java.util.UUID;

@Repository
public interface EntityRepository extends JpaRepository<Entity, UUID> {
}
