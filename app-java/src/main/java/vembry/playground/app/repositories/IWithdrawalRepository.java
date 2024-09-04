package vembry.playground.app.repositories;

import org.springframework.data.jpa.repository.JpaRepository;

import vembry.playground.app.models.Withdrawal;

public interface IWithdrawalRepository extends JpaRepository<Withdrawal, String> {
}
