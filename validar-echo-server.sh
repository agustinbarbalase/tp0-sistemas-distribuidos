#!/bin/bash
docker run --network=tp0_testing_net --rm 0ed463b26dae sh -c "echo 'Hola' | nc server 12345"