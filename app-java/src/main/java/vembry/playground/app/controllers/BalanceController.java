package vembry.playground.app.controllers;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import vembry.playground.app.models.Balance;
import vembry.playground.app.models.Deposit;
import vembry.playground.app.models.DepositRequest;
import vembry.playground.app.models.HttpResponse;
import vembry.playground.app.models.WithdrawRequest;
import vembry.playground.app.models.Withdrawal;
import vembry.playground.app.services.BalanceService;

@RestController
@RequestMapping("/balance")
public class BalanceController {

    @Autowired
    private BalanceService balanceService;

    @GetMapping("/{id}")
    public HttpResponse<Balance> get(@PathVariable String id) {
        Balance balance = this.balanceService.get(id);
        return new HttpResponse<Balance>(balance);
    }

    @PostMapping("/{id}/deposit")
    public HttpResponse<Deposit> deposit(@PathVariable String id, @RequestBody DepositRequest request) {
        request.setBalanceId(id);
        Deposit deposit = this.balanceService.deposit(request);
        return new HttpResponse<Deposit>(deposit);
    }

    @PostMapping("/{id}/withdraw")
    public HttpResponse<Withdrawal> withdraw(@PathVariable String id, @RequestBody WithdrawRequest request) {
        request.setBalanceId(id);
        Withdrawal withdrawal = this.balanceService.withdraw(request);
        return new HttpResponse<Withdrawal>(withdrawal);
    }
}
