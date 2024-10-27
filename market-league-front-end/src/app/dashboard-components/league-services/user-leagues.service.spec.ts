import { TestBed } from '@angular/core/testing';

import { UserLeaguesService } from './user-leagues.service';

describe('UserLeaguesService', () => {
  let service: UserLeaguesService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(UserLeaguesService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
