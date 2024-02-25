import { ExceptionBase } from '@base';

export const EMAIL_IN_USE = 'REGISTRATION.EMAIL_ALREADY_IN_USE';

/**
 * Used to indicate that provided email already in use in system
 *
 * @class EmailInUsedException
 * @extends {ExceptionBase}
 */
export class EmailInUsedException extends ExceptionBase {
  readonly code = EMAIL_IN_USE;
}