# Image Verification

The Spectromate container image is signed using [Sigstore's](https://sigstore.dev/) Cosign. The container image is signed using a cryptographic key pair that is private and stored internally. The public key is available in the official Spectro Cloud documentation repository at [**static/cosign.pub**](https://raw.githubusercontent.com/spectrocloud/librarium/master/static/cosign.pub). Use the public key to verify the authenticity of the container image. You can learn more about the container image signing process by reviewing the [Signing Containers](https://docs.sigstore.dev/signing/signing_with_containers) documentation page.


> [!NOTE]
> Cosign generates a key pair that uses the ECDSA-P256 algorithm for the signature and SHA256 for hashes. The keys are stored in PEM-encoded PKCS8 format.


Use the following command to verify the authenticity of the container image. Replace the image tag with the version you want to verify.

```shell
cosign verify --key https://raw.githubusercontent.com/spectrocloud/librarium/master/static/cosign.pub \
ghcr.io/spectrocloud/spectromate:v1.0.7
```

If the container image is valid, the following output is displayed. The example output is formatted using `jq` to improve readability.

```shell hideClipboard
Verification for ghcr.io/spectrocloud/spectromate:v1.0.7 --
The following checks were performed on each of these signatures:
  - The cosign claims were validated
  - Existence of the claims in the transparency log was verified offline
  - The signatures were verified against the specified public key
[
  {
    "critical": {
      "identity": {
        "docker-reference": "ghcr.io/spectrocloud/spectromate"
      },
      "image": {
        "docker-manifest-digest": "sha256:285a95a8594883b3748138460182142f5a1b74f80761e2fecb1b86d3c9b9d191"
      },
      "type": "cosign container image signature"
    },
    "optional": {
      "Bundle": {
        "SignedEntryTimestamp": "MEYCIQCZ6FZzNB5wA9+W/lF57jx0qTaszZhg5FxJiBmgIFxPVwIhANnoQQ5gqjr1h93LCq1Td8BohqrxxIvfrXTnT1tYR4i7",
        "Payload": {
          "body": "eyJhcGlWZXJzaW9uIjoiMC4wLjEiLCJraW5kIjoiaGFzaGVkcmVrb3JkIiwic3BlYyI6eyJkYXRhIjp7Imhhc2giOnsiYWxnb3JpdGhtIjoic2hhMjU2IiwidmFsdWUiOiI0MzU0MzFjNjY1Y2Y2ZGZjYzM0NzI2YWRkNjAzMDVjYjZlMzhlNjVkZmJlMWQ0NWU2ZGVkM2IzNzg3NTYwY2MxIn19LCJzaWduYXR1cmUiOnsiY29udGVudCI6Ik1FVUNJUUM0TFFxYVFDclhOc0VzdkI0ZE84bmtZSWg0L3o5UzdScGVEdUZnUDJwbDJ3SWdOdEJsNElDaHBmT3RnVDBlNW5QTmRMYWt4RTJHcnFFK0tjV1JXSGZPTnpnPSIsInB1YmxpY0tleSI6eyJjb250ZW50IjoiTFMwdExTMUNSVWRKVGlCUVZVSk1TVU1nUzBWWkxTMHRMUzBLVFVacmQwVjNXVWhMYjFwSmVtb3dRMEZSV1VsTGIxcEplbW93UkVGUlkwUlJaMEZGV1VoeVl6SlhTVVV6WVhCTFRHMWplR3hHUmtoNVZsRkRVVnBYYUFveUsyRnNOVmN2Vmsxc1VISXpkVFJGV2k5V0wwZFBRbTAySzFrNVowWXpWWE16ZEhkMVpWaFpaMlJaWlVadk5XODNRbFZ1TnpCTlVGQjNQVDBLTFMwdExTMUZUa1FnVUZWQ1RFbERJRXRGV1MwdExTMHRDZz09In19fX0=",
          "integratedTime": 1702758491,
          "logIndex": 57230483,
          "logID": "c0d23d6ad406973f9559f3ba2d1ca01f84147d8ffc5b8445c224f98b9591801d"
        }
      },
      "owner": "Spectro Cloud",
      "ref": "e597f70be238369ce4f0e5778492a155e23fec17",
      "repo": "spectrocloud/spectromate",
      "workflow": "Release"
    }
  }
]
```


> [!CAUTION]
> Do not use the container image if the authenticity cannot be verified. Verify you downloaded the correct public key and that the container image is from `ghcr.io/spectrocloud/spectromate`.

If the container image is not valid, an error is displayed. The following example shows an error when the container image is not valid.

```shell hideClipboard
cosign verify --key https://raw.githubusercontent.com/spectrocloud/librarium/master/static/cosign.pub \
ghcr.io/spectrocloud/spectromate:v1.0.66
```

```shell hideClipboard
Error: no matching signatures: error verifying bundle: comparing public key PEMs, expected -----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEheVfGYrVn2mIUQ4cxMJ6x09oXf82
zFEMG++p4q8Mf+y2gp7Ae4oUaXk6Q9V7aVjjltRVN6SQcoSASxf2H2EpgA==
-----END PUBLIC KEY-----
, got -----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEYHrc2WIE3apKLmcxlFFHyVQCQZWh
2+al5W/VMlPr3u4EZ/V/GOBm6+Y9gF3Us3twueXYgdYeFo5o7BUn70MPPw==
-----END PUBLIC KEY-----

main.go:69: error during command execution: no matching signatures: error verifying bundle: comparing public key PEMs, expected -----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEheVfGYrVn2mIUQ4cxMJ6x09oXf82
zFEMG++p4q8Mf+y2gp7Ae4oUaXk6Q9V7aVjjltRVN6SQcoSASxf2H2EpgA==
-----END PUBLIC KEY-----
, got -----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEYHrc2WIE3apKLmcxlFFHyVQCQZWh
2+al5W/VMlPr3u4EZ/V/GOBm6+Y9gF3Us3twueXYgdYeFo5o7BUn70MPPw==
-----END PUBLIC KEY-----
```
