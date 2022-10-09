import {PostConfirmationConfirmSignUpTriggerEvent, PreSignUpTriggerEvent} from "aws-lambda";

exports.handler = async (event: PostConfirmationConfirmSignUpTriggerEvent) => {
    console.log(event)
    return event;
}