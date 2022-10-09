import {PreSignUpTriggerEvent} from 'aws-lambda';

exports.handler = async (event: PreSignUpTriggerEvent) => {
    event.response.autoConfirmUser = true;
    event.response.autoVerifyPhone = true;
    console.info(JSON.stringify(event));
    return event;
}