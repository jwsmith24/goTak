# Behavior-test reference

A durable test:

- names a capability in domain language;
- exercises a public interface at an agreed seam;
- derives its expectation from a specification, worked example, invariant, or known literal independent of the implementation;
- observes externally meaningful output or state through that interface;
- survives behavior-preserving internal refactors.

Treat a test that mocks internal collaborators, calls private methods, asserts internal call order, queries around the interface, or recreates the implementation in its expectation as a design warning.