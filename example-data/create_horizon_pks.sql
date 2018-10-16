--need to update the horizon db to have primary keys
alter table history_accounts ADD PRIMARY KEY (id);
alter table history_effects ADD PRIMARY KEY (history_account_id, history_operation_id, "order");
alter table history_ledgers ADD PRIMARY KEY (id);
alter table history_operations ADD PRIMARY KEY (id);
alter table history_trades ADD PRIMARY KEY (history_operation_id, "order");
alter table history_transactions ADD PRIMARY KEY (id);