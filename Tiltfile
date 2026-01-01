# Load the restart_process extension
load('ext://restart_process', 'docker_build_with_restart')

# Set the namespace
k8s_namespace('ride-sharing-dev')

# -------------------------------------------------------------------------
# Helm Chart Integration
# -------------------------------------------------------------------------
# Path to the Helm chart relative to this Tiltfile
chart_path = '../ride-sharing-platform/charts/ride-sharing'

# Load Deployment Manifests from Helm
# We combine values-dev.yaml (config) and secrets-dev.yaml (secrets)
# Load Deployment Manifests from Helm using local command (more reliable)
# We execute the exact same command that works in your terminal
# Load Deployment Manifests from Helm using local command
# Debug: Output to file so we can inspect it
manifests = local(
    'helm template ride-sharing ../ride-sharing-platform/charts/ride-sharing ' +
    '-f ../ride-sharing-platform/charts/ride-sharing/values-dev.yaml ' +
    '-f ../ride-sharing-platform/charts/ride-sharing/values-secrets-dev.yaml ' +
    '--namespace ride-sharing-dev'
)

# Apply the generated YAML to the cluster
k8s_yaml(manifests)

# -------------------------------------------------------------------------
# Infrastructure (No Build Required)
# -------------------------------------------------------------------------

# RabbitMQ
k8s_resource('rabbitmq', 
    port_forwards=['5672', '15672'], 
    labels='tooling'
)

# Jaeger
k8s_resource('jaeger', 
    port_forwards=['16686:16686', '14268:14268'], 
    labels='tooling'
)

# ConfigMaps/Secrets are managed by Helm, no need to specific resource definition
# unless we want to trigger updates on them specifically.
# They are auto-deployed via k8s_yaml(manifests)

# -------------------------------------------------------------------------
# Microservices (Build & Deploy)
# -------------------------------------------------------------------------

# --- API Gateway ---
gateway_compile_cmd = 'CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/api-gateway ./services/api-gateway'
if os.name == 'nt':
  gateway_compile_cmd = './infra/development/docker/api-gateway-build.bat'

local_resource(
  'api-gateway-compile',
  gateway_compile_cmd,
  deps=['./services/api-gateway', './shared'], 
  labels="compiles"
)

docker_build_with_restart(
  'ride-sharing/api-gateway', # Matches image in values.yaml
  '.',
  entrypoint=['/app/build/api-gateway'],
  dockerfile='./infra/development/docker/api-gateway.Dockerfile',
  only=['./build/api-gateway', './shared'],
  live_update=[
    sync('./build', '/app/build'),
    sync('./shared', '/app/shared'),
  ],
)

k8s_resource('api-gateway', 
    port_forwards=8081,
    resource_deps=['api-gateway-compile', 'rabbitmq'], 
    labels="services"
)


# --- Driver Service ---
driver_compile_cmd = 'CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/driver-service ./services/driver-service'
if os.name == 'nt':
 driver_compile_cmd = './infra/development/docker/driver-build.bat'

local_resource(
  'driver-service-compile',
  driver_compile_cmd,
  deps=['./services/driver-service', './shared'], 
  labels="compiles"
)

docker_build_with_restart(
  'ride-sharing/driver-service',
  '.',
  entrypoint=['/app/build/driver-service'],
  dockerfile='./infra/development/docker/driver-service.Dockerfile',
  only=['./build/driver-service', './shared'],
  live_update=[
    sync('./build', '/app/build'),
    sync('./shared', '/app/shared'),
  ],
)

k8s_resource('driver-service', 
    resource_deps=['driver-service-compile', 'rabbitmq'], 
    labels="services"
)


# --- Payment Service ---
payment_compile_cmd = 'CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/payment-service ./services/payment-service/cmd/main.go'
if os.name == 'nt':
  payment_compile_cmd = './infra/development/docker/payment-build.bat'

local_resource(
  'payment-service-compile',
  payment_compile_cmd,
  deps=['./services/payment-service', './shared'], 
  labels="compiles"
)

docker_build_with_restart(
  'ride-sharing/payment-service',
  '.',
  entrypoint=['/app/build/payment-service'],
  dockerfile='./infra/development/docker/payment-service.Dockerfile',
  only=['./build/payment-service', './shared'],
  live_update=[
    sync('./build', '/app/build'),
    sync('./shared', '/app/shared'),
  ],
)

k8s_resource('payment-service', 
    resource_deps=['payment-service-compile', 'rabbitmq'], 
    labels="services"
)


# --- Trip Service ---
trip_compile_cmd = 'CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/trip-service ./services/trip-service/cmd/main.go'
if os.name == 'nt':
 trip_compile_cmd = './infra/development/docker/trip-build.bat'

local_resource(
  'trip-service-compile',
  trip_compile_cmd,
  deps=['./services/trip-service', './shared'], 
  labels="compiles"
)

docker_build_with_restart(
  'ride-sharing/trip-service',
  '.',
  entrypoint=['/app/build/trip-service'],
  dockerfile='./infra/development/docker/trip-service.Dockerfile',
  only=['./build/trip-service', './shared'],
  live_update=[
    sync('./build', '/app/build'),
    sync('./shared', '/app/shared'),
  ],
)

k8s_resource('trip-service', 
    resource_deps=['trip-service-compile', 'rabbitmq'], 
    labels="services"
)


# --- Web Frontend ---
docker_build(
  'ride-sharing/web',
  '.',
  dockerfile='./infra/development/docker/web.Dockerfile',
)

k8s_resource('web', 
    port_forwards=3000, 
    labels="frontend"
)