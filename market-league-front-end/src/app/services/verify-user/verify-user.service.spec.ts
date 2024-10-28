import { TestBed } from '@angular/core/testing';

import { VerifyUserService } from './verify-user.service';

describe('VerifyUserService', () => {
  let service: VerifyUserService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(VerifyUserService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
