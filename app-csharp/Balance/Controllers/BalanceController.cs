using Microsoft.AspNetCore.Mvc;

namespace App.Balance.Controllers
{
    [Route("balance")]
    [ApiController]
    public class BalanceController : ControllerBase
    {
        [HttpPost("open")]
        public ActionResult<string> OpenBalance()
        {
            return "opening balance";
        }

        [HttpGet("{balanceId}")]
        public ActionResult<string> GetBalance(string balanceId)
        {
            return balanceId;
        }

        [HttpPost("{balanceId}/deposit")]
        public ActionResult<string> Deposit(string balanceId)
        {
            return balanceId;
        }

        [HttpPost("{balanceId}/withdraw")]
        public ActionResult<string> Withdraw(string balanceId)
        {
            return balanceId;
        }

        [HttpPost("{balanceIdFrom}/transfer/{balanceIdTo}")]
        public ActionResult<string> Transfer(string balanceIdFrom, string balanceIdTo)
        {
            return balanceIdFrom;
        }
    
    }
}