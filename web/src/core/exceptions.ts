import { ExceptionBase } from '@base'

export const ARGUMENT_INVALID = 'GENERIC.ARGUMENT_INVALID';
export const ARGUMENT_NOT_PROVIDED = 'GENERIC.ARGUMENT_NOT_PROVIDED';
export const UNEXPECTED_ERROR = 'GENERIC.UNEXPECTED_ERROR';

/**
 * Used to indicate that an incorrect argument was provided to a method/function/class constructor
 *
 * @class ArgumentInvalidException
 * @extends {ExceptionBase}
 */
export class ArgumentInvalidException extends ExceptionBase {
  readonly code = ARGUMENT_INVALID;
}

/**
 * Used to indicate that an argument was not provided (is empty object/array, null of undefined).
 *
 * @class ArgumentNotProvidedException
 * @extends {ExceptionBase}
 */
export class ArgumentNotProvidedException extends ExceptionBase {
  readonly code = ARGUMENT_NOT_PROVIDED;
}

/**
 * Used to indicate an unexpected error that does not fall under all other errors
 *
 * @class UnexpectedErrorException
 * @extends {ExceptionBase}
 */
export class UnexpectedErrorException extends ExceptionBase {
  static readonly message = 'Unexpected error';

  constructor(message = UnexpectedErrorException.message) {
    super(message);
  }

  readonly code = UNEXPECTED_ERROR;
}
