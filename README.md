# Traefik Plugin Response Cache Control

This plugin allows you to set or override Cache-Control headers in HTTP responses.

## Configuration

```yaml
http:
  middlewares:
    response-cache-control:
      plugin:
        response-cache-control:
          value: "public, max-age=3600"
          override: true
          excludedStatusCodes: ["400-499", "500"]
```

### Configuration Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `value` | String | `"public, max-age=3600"` | The Cache-Control header value to be set |
| `override` | Boolean | `true` | Whether to override existing Cache-Control headers |
| `excludedStatusCodes` | []String | `[]` | Status codes or ranges to exclude from Cache-Control header setting. Format can be single codes (`"404"`) or ranges (`"400-499"`) |

## Example Usage

### Basic Usage with Default Values

```yaml
http:
  middlewares:
    my-cache-control:
      plugin:
        response-cache-control: {}
```

### Custom Cache Configuration

```yaml
http:
  middlewares:
    static-cache:
      plugin:
        response-cache-control:
          value: "public, max-age=86400"
          override: true
          excludedStatusCodes: ["404", "500-599"]
```

### Using with a Traefik Router

```yaml
http:
  routers:
    my-router:
      rule: "Host(`example.com`)"
      middlewares:
        - "response-cache-control"
      service: "my-service"
```

## Development Notes

When developing this plugin, please note that the package name must be `traefik_plugin_response_cache_control` to match the module path's last component. This is required for proper integration with Traefik's Yaegi interpreter.

```go
// Correct package declaration
package traefik_plugin_response_cache_control
```
