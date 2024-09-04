package vembry.playground.app.repositories;

import org.springframework.data.jpa.repository.JpaRepository;

import vembry.playground.app.models.Deposit;

public interface IDepositRepository extends JpaRepository<Deposit, String> {
}
