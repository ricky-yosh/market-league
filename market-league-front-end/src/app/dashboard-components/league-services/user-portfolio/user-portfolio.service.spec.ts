import { TestBed } from '@angular/core/testing';

import { UserPortfolioService } from './user-portfolio.service';

describe('UserPortfolioService', () => {
  let service: UserPortfolioService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(UserPortfolioService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
