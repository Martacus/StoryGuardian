package com.guardian.gaia.repository;

import com.guardian.custodian.Story;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import java.util.UUID;

@Repository
public interface StoryRepository extends JpaRepository<Story, UUID> {
}
