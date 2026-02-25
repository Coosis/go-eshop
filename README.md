# simple online store backend
## components:
- Catalog
- Stock
- Seckill
- Order
- Cart

# load testing with k6:
```
=== summary_b250_c31.json ===
http_reqs_rate=901.0790355945727
http_req_failed=0.0015531047320755831
p95=10.61890525
p99=n/a
=== summary_b500_c63.json ===
http_reqs_rate=1788.4699056295115
http_req_failed=0.010818903631083482
p95=18.851899199999988
p99=n/a
```

Then my load generator ran out of memory and was killed by the os.
For logs/statistics/graphs, you can refer to: 
[run_20260216](https://pub-b15d04bbf0384cf3934c0584b84110f1.r2.dev/run_20260216.tar.gz).
