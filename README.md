# cs686-blockchain-p5

Project 5 Proposal
-------------------------

What -> 

Lending system on blockchain


Why ->

Borrowing and Lending industry is very inefficient. Lot of scope in making it efficient.

On blockchain, lenders can remain anonymous and can even invest in small amounts, without the hassle of banking system and country politics. Lenders get interest on the investment.

Borrowers can get money from different sources, without relying on traditional lending systems, that too with quite low interest rate.


How ->

Borrowers and Lenders can register themselves on network. 

Miners will get tx fees, where tx is defined as 

Lenders giving money and Corresponding amount being credited on borrowers side.


List of functionalities ->

Borrowers and Lenders registers
Miners/Validators/etc - will put the tx on the blockchain
Borrowers can request their required amount.
Many Lenders can contribute to the 1 borrowers requirement.
Wait until Lenders have contributed more than or equal to borrowers requirement.
If in specified time period the borrower have enough lending promised, then release the fund to borrower.
After certain time(set by borrower as part of request), the Lenders are paid interest proportional to lend amount. / (or not if the borrower does not make profit.) - in that case will the borrower have to part with some share of the company.


Block consist of txs. Txs contain information either - borrowing requirement Or lenders promising money to a requirement (at that point money for a borrower is stored in a separate account), when total of txs exceed the required amount, the amount is released to the borrower. Borrower can then push the expense on blockchain to investors and borrower can maintain trust. After certain (pre decided time), interest will be given back to the lenders.

Measure of success ->

schedule - 

18th Apr - 24th Apr	- Registration system (Allow borrower and Lender to register, registration of company, and generating holding account for each company) 

25th Apr - 1st May	- Each party can interact with system, borrowers can create requirement, Lenders can start lending against that requirement, the lend amount should go to holding account, until amount becomes greater than requirement, then lended amount gets moved to borrower’s company account

2nd May - 8th May - TBD

8th May - 15th May	- TBD
