package vembry.playground.app.repositories;

import org.springframework.data.jpa.repository.JpaRepository;

import vembry.playground.app.models.Balance;

public interface IBalanceRepository extends JpaRepository<Balance, String> {
}
