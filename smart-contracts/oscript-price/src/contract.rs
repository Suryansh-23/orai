use crate::error::ContractError;
use crate::msg::{HandleMsg, InitMsg, QueryMsg};
use crate::state::{config, config_read, State};
use cosmwasm_std::{
    to_binary, Api, Binary, Env, Extern, HandleResponse, InitResponse, MessageInfo, Querier,
    StdResult, Storage,
};

// Note, you can use StdResult in some functions where you do not
// make use of the custom errors
pub fn init<S: Storage, A: Api, Q: Querier>(
    deps: &mut Extern<S, A, Q>,
    _env: Env,
    info: MessageInfo,
    msg: InitMsg,
) -> StdResult<InitResponse> {
    let state = State {
        ai_data_source: msg.ai_data_source,
        testcase: msg.testcase,
        owner: deps.api.canonical_address(&info.sender)?,
    };
    config(&mut deps.storage).save(&state)?;

    Ok(InitResponse::default())
}

// And declare a custom Error variant for the ones where you will want to make use of it
pub fn handle<S: Storage, A: Api, Q: Querier>(
    deps: &mut Extern<S, A, Q>,
    _env: Env,
    info: MessageInfo,
    msg: HandleMsg,
) -> Result<HandleResponse, ContractError> {
    match msg {
        HandleMsg::UpdateDatasource { name } => try_update_datasource(deps, info, name),
        HandleMsg::UpdateTestcase { name } => try_update_testcase(deps, info, name),
    }
}

pub fn try_update_datasource<S: Storage, A: Api, Q: Querier>(
    deps: &mut Extern<S, A, Q>,
    info: MessageInfo,
    name: String,
) -> Result<HandleResponse, ContractError> {
    let api = &deps.api;
    config(&mut deps.storage).update(|mut state| -> Result<_, ContractError> {
        if api.canonical_address(&info.sender)? != state.owner {
            return Err(ContractError::Unauthorized {});
        }
        state.ai_data_source = name;
        Ok(state)
    })?;
    Ok(HandleResponse::default())
}

pub fn try_update_testcase<S: Storage, A: Api, Q: Querier>(
    deps: &mut Extern<S, A, Q>,
    info: MessageInfo,
    name: String,
) -> Result<HandleResponse, ContractError> {
    let api = &deps.api;
    config(&mut deps.storage).update(|mut state| -> Result<_, ContractError> {
        if api.canonical_address(&info.sender)? != state.owner {
            return Err(ContractError::Unauthorized {});
        }
        state.testcase = name;
        Ok(state)
    })?;
    Ok(HandleResponse::default())
}

pub fn query<S: Storage, A: Api, Q: Querier>(
    deps: &Extern<S, A, Q>,
    _env: Env,
    msg: QueryMsg,
) -> StdResult<Binary> {
    match msg {
        QueryMsg::GetDatasource {} => to_binary(&query_datasource(deps)?),
        QueryMsg::GetTestcase {} => to_binary(&query_testcase(deps)?),
        QueryMsg::Aggregate { results } => to_binary(&query_aggregation(deps, results)?),
    }
}

fn query_datasource<S: Storage, A: Api, Q: Querier>(deps: &Extern<S, A, Q>) -> StdResult<String> {
    let state = config_read(&deps.storage).load()?;
    Ok(state.ai_data_source)
}

fn query_testcase<S: Storage, A: Api, Q: Querier>(deps: &Extern<S, A, Q>) -> StdResult<String> {
    let state = config_read(&deps.storage).load()?;
    Ok(state.testcase)
}

fn query_aggregation<S: Storage, A: Api, Q: Querier>(
    _deps: &Extern<S, A, Q>,
    results: Vec<String>,
) -> StdResult<String> {
    let mut sum: i32 = 0;
    let mut floating_sum: i32 = 0;
    let mut count = 0;
    for input in results {
        // get first item from iterator
        let mut iter = input.split('.');
        let first = iter.next();
        let last = iter.next();
        // will panic instead for forward error with ?
        let number: i32 = first.unwrap().parse().unwrap();
        let mut floating: i32 = 0;
        if last.is_some() {
            let mut last_part = last.unwrap().to_owned();
            if last_part.len() < 2 {
                last_part.push_str("0");
            } else if last_part.len() > 2 {
                last_part = last_part[..2].to_string();
            }
            floating = last_part.parse().unwrap();
        }
        sum += number;
        floating_sum += floating;
        count += 1;
    }

    sum = sum / count;
    floating_sum = floating_sum / count;
    let final_result = format!("{}.{}", sum, floating_sum);

    Ok(final_result)
}

#[cfg(test)]
mod tests {
    use super::*;
    use cosmwasm_std::testing::{mock_dependencies, mock_env, mock_info};
    use cosmwasm_std::{coins, from_binary};

    #[test]
    fn proper_initialization() {
        let mut deps = mock_dependencies(&[]);

        let msg = InitMsg {
            ai_data_source: "datasource_eth".to_string(),
            testcase: "testcase_price".to_string(),
        };
        let info = mock_info("creator", &coins(1000, "earth"));

        // we can just call .unwrap() to assert this was a success
        let res = init(&mut deps, mock_env(), info, msg).unwrap();
        assert_eq!(0, res.messages.len());

        // it worked, let's query the state
        let res = query(&deps, mock_env(), QueryMsg::GetDatasource {}).unwrap();
        let value: String = from_binary(&res).unwrap();
        assert_eq!("datasource_eth", value);
    }

    #[test]
    fn update_datasource() {
        let mut deps = mock_dependencies(&coins(2, "token"));

        let msg = InitMsg {
            ai_data_source: "datasource_eth".to_string(),
            testcase: "testcase_price".to_string(),
        };
        let info = mock_info("creator", &coins(2, "token"));
        let _res = init(&mut deps, mock_env(), info, msg).unwrap();

        // beneficiary can release it
        let info = mock_info("creator", &coins(2, "token"));
        let msg = HandleMsg::UpdateDatasource {
            name: "datasource_btc".to_string(),
        };
        let _res = handle(&mut deps, mock_env(), info, msg).unwrap();

        // should increase counter by 1
        let res = query(&deps, mock_env(), QueryMsg::GetDatasource {}).unwrap();
        let value: String = from_binary(&res).unwrap();
        assert_eq!("datasource_btc", value);
    }
}
