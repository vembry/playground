package vembry.playground.app.services;

import org.springframework.amqp.rabbit.core.RabbitTemplate;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import vembry.playground.app.models.Balance;
import vembry.playground.app.models.Deposit;
import vembry.playground.app.models.DepositRequest;
import vembry.playground.app.models.WithdrawRequest;
import vembry.playground.app.models.Withdrawal;
import vembry.playground.app.repositories.IBalanceRepository;
import vembry.playground.app.repositories.IDepositRepository;
import vembry.playground.app.repositories.IWithdrawalRepository;

@Service
public class BalanceService {

    @Autowired
    private IBalanceRepository balanceRepository;

    @Autowired
    private IDepositRepository depositRepository;

    @Autowired
    private IWithdrawalRepository withdrawalRepository;

    @Autowired
    private RabbitTemplate rabbitTemplate;

    public Balance get(String id) {
        Balance balance = this.balanceRepository.findById(id).orElse(null);
        return balance;
    }

    public Deposit deposit(DepositRequest request) {
        Deposit deposit = new Deposit(request.getBalanceId(), request.getAmount());
        deposit = this.depositRepository.save(deposit);
        this.rabbitTemplate.convertAndSend("deposit", deposit.getId());
        return deposit;
    }

    public Withdrawal withdraw(WithdrawRequest request) {
        Withdrawal withdrawal = new Withdrawal(request.getBalanceId(), request.getAmount());
        withdrawal = this.withdrawalRepository.save(withdrawal);
        this.rabbitTemplate.convertAndSend("withdraw", withdrawal.getId());
        return withdrawal;
    }
}
